#! /bin/bash

cleanup() {
    ./.semaphore/collect-artifacts
    cp -a /src/calico/felix/report/. /src/calico/felix/artifacts
    ./.semaphore/publish-artifacts-argoci
}

checkExitCode() {
    if [ $1 -ne 0 ]; then
        export CI_EXIT_CODE=1
        cleanup
        return 1
    fi
}

cd felix

docker load -i /src/dockerimages/calico-felix.tar
docker tag calico/felix:latest-amd64 felix:latest-amd64
docker load -i /src/dockerimages/felixtest-typha.tar
docker tag felix-test/typha:latest-amd64 typha:latest-amd64

make check-wireguard
checkExitCode $? || return 0
../.argoci/scripts/run-and-monitor fv-${JOB_INDEX}.log make fv-no-prereqs FV_BATCHES_TO_RUN="${JOB_INDEX}" FV_NUM_BATCHES=${JOB_COUNT}
checkExitCode $? || return 0

cleanup
