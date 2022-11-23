#! /bin/bash

export GOOGLE_APPLICATION_CREDENTIALS=$HOME/secrets/secret.google-service-account-key.json
export SHORT_WORKFLOW_ID=$(echo {{workflow.uid}} | sha256sum | cut -c -8)
export ZONE=europe-west3-c
export VM_PREFIX=argoci-${REPO_NAME}-${SHORT_WORKFLOW_ID}-
echo VM_PREFIX=${VM_PREFIX}
export NUM_FV_BATCHES=8

cd felix
mkdir artifacts

# Create initial VMs
./.semaphore/create-test-vms ${VM_PREFIX}
status=$?
if [ $status -ne 0 ]; then
    # VM Creation failed - Exit script now and global_epilogue will upload the artifacts
    # Set CI_EXIT_CODE so we know to mark the failure
    export CI_EXIT_CODE=1
    ./.semaphore/publish-artifacts-argoci
    return
fi

./.semaphore/run-tests-on-vms ${VM_PREFIX}
export CI_EXIT_CODE=$?

./.semaphore/collect-artifacts-from-vms ${VM_PREFIX}
./.semaphore/publish-artifacts-argoci
./.semaphore/clean-up-vms ${VM_PREFIX}