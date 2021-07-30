import org.kie.jenkins.jobdsl.templates.KogitoJobTemplate
import org.kie.jenkins.jobdsl.KogitoConstants
import org.kie.jenkins.jobdsl.Utils

boolean isMainBranch() {
    return "${GIT_BRANCH}" == "${GIT_MAIN_BRANCH}"
}

def getDefaultJobParams() {
    return [
        job: [
            name: 'rhpam-kogito-operator'
        ],
        git: [
            author: "${GIT_AUTHOR_NAME}",
            branch: "${GIT_BRANCH}",
            repository: 'rhpam-kogito-operator',
            credentials: "${GIT_AUTHOR_CREDENTIALS_ID}",
            token_credentials: "${GIT_AUTHOR_TOKEN_CREDENTIALS_ID}"
        ]
    ]
}

def getJobParams(String jobName, String jobFolder, String jenkinsfileName, String jobDescription = '') {
    def jobParams = getDefaultJobParams()
    jobParams.job.name = jobName
    jobParams.job.folder = jobFolder
    jobParams.jenkinsfile = jenkinsfileName
    if (jobDescription) {
        jobParams.job.description = jobDescription
    }
    return jobParams
}

def nightlyBranchFolder = "${KogitoConstants.KOGITO_DSL_NIGHTLY_FOLDER}/${JOB_BRANCH_FOLDER}"

// if (isMainBranch()) {
    // PR job is disabled for now as handled by another Jenkins
    // folder(KogitoConstants.KOGITO_DSL_PULLREQUEST_FOLDER)

    // setupPrJob(KogitoConstants.KOGITO_DSL_PULLREQUEST_FOLDER)
// }

setupSyncJob(nightlyBranchFolder)

setupProdUpdateVersionJob("${KogitoConstants.KOGITO_DSL_TOOLS_FOLDER}/${JOB_BRANCH_FOLDER}")

/////////////////////////////////////////////////////////////////
// Methods
/////////////////////////////////////////////////////////////////

void setupPrJob(String jobFolder) {
    def jobParams = getDefaultJobParams()
    jobParams.job.folder = jobFolder
    KogitoJobTemplate.createPRJob(this, jobParams)
}


void setupSyncJob(String jobFolder) {
    def jobParams = getJobParams('rhpam-kogito-operator-sync', jobFolder, 'Jenkinsfile.upstream-operator-sync', 'RHPAM Kogito Operator synchronizing with Kogito operator')
    jobParams.triggers = [ cron : '@midnight' ]
    KogitoJobTemplate.createPipelineJob(this, jobParams).with {
        parameters {
            stringParam('DISPLAY_NAME', '', 'Setup a specific build display name')

            stringParam('BUILD_BRANCH_NAME', "${GIT_BRANCH}", 'Set the Git branch to checkout')

            // Build&Test information
            stringParam('BDD_TEST_TAGS', '@rhpam', 'Execute only a subset of BDD tests')
        }

        environmentVariables {
            env('JENKINS_EMAIL_CREDS_ID', "${JENKINS_EMAIL_CREDS_ID}")

            env('REPO_NAME', 'rhpam-kogito-operator')
            env('OPERATOR_IMAGE_NAME', 'rhpam-kogito-rhel8-operator')
            env('CONTAINER_ENGINE', 'podman')
            env('CONTAINER_TLS_OPTIONS', '--tls-verify=false')
            env('OPENSHIFT_API_KEY', 'OPENSHIFT_API')
            env('OPENSHIFT_CREDS_KEY', 'OPENSHIFT_CREDS')
            env('OPENSHIFT_REGISTRY_KEY', 'OPENSHIFT_REGISTRY')

            env('GIT_AUTHOR', "${GIT_AUTHOR_NAME}")
        }
    }
}

void setupProdUpdateVersionJob(String jobFolder) {
    KogitoJobTemplate.createPipelineJob(this, getJobParams('rhpam-kogito-operator-update-prod-version', jobFolder, 'Jenkinsfile.update-prod-version', 'Update prod version for RHPAM Kogito Operator')).with {
        parameters {
            stringParam('PROD_PROJECT_VERSION', '', 'Which version to set ?')
            stringParam('PROD_BUNDLE_SUFFIX', '', '(Optional) Bundle suffix to apply to the version ? Default is value `1`.')
            stringParam('PROD_REPLACES_VERSION', '', '(Optional) Which version does it replaces ? If not given, no replacement will be done.')
        }

        environmentVariables {
            env('REPO_NAME', 'kogito-images')

            env('GIT_AUTHOR', "${GIT_AUTHOR_NAME}")
            env('BUILD_BRANCH_NAME', "${GIT_BRANCH}")

            env('AUTHOR_CREDS_ID', "${GIT_AUTHOR_CREDENTIALS_ID}")
            env('GITHUB_TOKEN_CREDS_ID', "${GIT_AUTHOR_TOKEN_CREDENTIALS_ID}")
            env('GIT_AUTHOR_BOT', "${GIT_BOT_AUTHOR_NAME}")
            env('BOT_CREDENTIALS_ID', "${GIT_BOT_AUTHOR_CREDENTIALS_ID}")
        }
    }
}
