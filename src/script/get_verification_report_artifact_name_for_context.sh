#!/bin/bash

# This script attempt to generate a verification report artifact name (to which the report is uploaded in ADO) 
# based on the given arguments.
#
# arg1: Ready for [use|production]
#
# Note that "production" in this context means that the verification report indicates readyness to deploy to the production environment.
# Likewise, "use" in this context means that the verification report indicates readyness to release what is deployed to the production environment for use.
#
# Example usage:
# ./get_verification_report_artifact_name_for_context.sh "production"
#    > verification_report_validation
# ./get_verification_report_artifact_name_for_context.sh "use"
#    > verification_report_production


#################
### FUNCTIONS ###
#################

# Takes parameter: "$READY_FOR"
generate_verification_report_artifact_name () {
    # Determine execution context
    if [ "$READY_FOR" == "production" ]; then
        # Context: Generated for validation environment
        # Generate the name of the verification report artifact for production
        vran=verification_report_validation
    elif [ "$READY_FOR" == "use" ]; then
        # Context: Generated for production environment
        # Generate the name of the verification report artifact for validation
        vran=verification_report_production
    else
        # Nothing could be determined, so exit with error
        vran=""
        echo 'ERROR: Could not determine execution context'
        exit 1
    fi
}



#############
### START ###
#############

# Get the parameter values
READY_FOR=$1

# Guard block: 
# Check if the parameter values are not empty
if [[ -z $READY_FOR ]]; then
    echo 'One or more parameters are undefined'
    exit 1
fi

# Generate
generate_verification_report_artifact_name "$READY_FOR"
echo $vran