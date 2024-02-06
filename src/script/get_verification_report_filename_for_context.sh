#!/bin/bash

# This script attempt to generate a verification report name based on the given arguments.
#
# arg1: Environment name
# arg2: Build ID
# arg3: Ready for [use|production]
#
# Note that "production" in this context means that the verification report indicates readyness to deploy to the production environment.
# Likewise, "use" in this context means that the verification report indicates readyness to release what is deployed to the production environment for use.
#
# Example usage:
# ./get_verification_report_filename_for_context.sh "ramone-service1-val-eu-central1" "617829" "use"
#    > VerificationReport_production_617829_ramone_service1_val_eu_central1
# ./get_verification_report_filename_for_context.sh "ramone-service1-val-eu-central1" "617829" "production"
#    > VerificationReport_validation_617829_ramone_service1_val_eu_central1


#################
### FUNCTIONS ###
#################

# Takes parameters: "$ENVIRONMENT_NAME" "$BUILD_ID" "$READY_FOR"
generate_verification_report_name () {
    # Determine execution context
    if [ "$READY_FOR" == "production" ]; then
        # Context: Generated for validation environment
        # Generate the name of the verification report for artifact in validation
        vrn=VerificationReport_validation_${BUILD_ID}_${ENVIRONMENT_NAME}
    elif [ "$READY_FOR" == "use" ]; then
        # Context: Generated for production environment
        # Generate the name of the verification report for artifact in production
        vrn=VerificationReport_production_${BUILD_ID}_${ENVIRONMENT_NAME}
    else
        # Nothing could be determined, so exit with error
        vrn=""
        echo 'ERROR: Could not determine execution context'
        exit 1
    fi
}



#############
### START ###
#############

# Get the parameter values
ENVIRONMENT_NAME=$1
BUILD_ID=$2
READY_FOR=$3

# Guard block: 
# Check if the parameter values are not empty
if [[ -z $BUILD_ID || -z $READY_FOR || -z $ENVIRONMENT_NAME ]]; then
    echo 'One or more parameters are undefined'
    exit 1
fi

# Preparation block:
# Replace all blanks
ENVIRONMENT_NAME=${ENVIRONMENT_NAME// /_}
# Replace all dashes
ENVIRONMENT_NAME=${ENVIRONMENT_NAME//-/_}
# To lowercase
ENVIRONMENT_NAME=$(echo "$ENVIRONMENT_NAME" | tr '[:upper:]' '[:lower:]')

# Generate
generate_verification_report_name "$ENVIRONMENT_NAME" "$BUILD_ID" "$READY_FOR"
echo $vrn