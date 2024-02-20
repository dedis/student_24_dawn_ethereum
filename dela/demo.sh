#!/usr/bin/env bash

set -e

export TEMPDIR=$(mktemp -d /tmp/dkgcli.XXXXXXXXXXXXX)
rm_tempdir () {
 rm -rf "$TEMPDIR"
}
trap rm_tempdir EXIT

n=16
t=9

for i in $(seq $n); do
	tmux new-window -d "LLVL=info dkgcli --config $TEMPDIR/node$i start --listen tcp://127.0.0.1:$((2000+i)); read"
done

sleep 3

# Exchange certificates
for i in $(seq 2 $n); do
	dkgcli --config $TEMPDIR/node$i minogrpc join --address //127.0.0.1:2001 $(dkgcli --config $TEMPDIR/node1 minogrpc token)
done

# Initialize DKG on each node. Do that in a 4th session.
for i in $(seq $n); do
	dkgcli --config $TEMPDIR/node$i dkg listen
done

# Do the setup in one of the node:
cmd=(dkgcli --config $TEMPDIR/node1 dkg setup --threshold $t) 
for i in $(seq $n); do
    cmd+=(--authority $(cat $TEMPDIR/node$i/dkgauthority))
done
"${cmd[@]}"

verbs=(buy sell)
animals=(penguin lion monkey cat dog kangaroo)
formats=(jpeg png gif)
coins=(poptokens byzcoins efranks)

function pick {
	shift $((RANDOM%$#))
	echo "$1"
}


echo ✅ generated committee pubkey: $(dkgcli --config $TEMPDIR/node1 dkg get-public-key)

for ((i = 0; i < 12; i++)); do
	sleep 1
	label=$(printf "%02x" $i)
	echo ⌚ block $i decryption key:  $(dkgcli --config $TEMPDIR/node1 dkg sign -message $label)
	case $i in
		( 5 )
			msg="$(pick ${verbs[@]}) a $(pick ${animals[@]}) $(pick ${formats[@]}) in exchange for $(pick ${coins[@]})"
			ct=$(dkgcli --config $TEMPDIR/node1 dkg encrypt -label 0a -message $(xxd -p -c0 <<<$msg))
			echo 🔐 encrypted message for block 10: $ct
			;;
		( 10 )
			echo 🔓 decrypted message: $(dkgcli --config $TEMPDIR/node1 dkg decrypt -label $label -ciphertext $ct)
			echo 🔓 decrypted message: $(dkgcli --config $TEMPDIR/node1 dkg decrypt -label $label -ciphertext $ct | xxd -r -p)
			;;
	esac
done

echo ⚠️ NOT FINANCIAL ADVICE⚠️

read
