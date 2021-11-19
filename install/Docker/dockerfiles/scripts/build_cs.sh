LIST=`ls | grep -E '^svc-'`
echo $LIST
for i in $LIST; do
    cd $i
    if [[ "$i" == "svc-composition-service" ]]; then
        /bin/bash build.sh
    fi
done