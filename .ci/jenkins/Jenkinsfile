@Library('jenkins-pipeline-shared-libraries')_

def changeAuthor = env.ghprbPullAuthorLogin ?: CHANGE_AUTHOR
def changeBranch = env.ghprbSourceBranch ?: CHANGE_BRANCH
def changeTarget = env.ghprbTargetBranch ?: CHANGE_TARGET

pipeline {
    agent { label 'operator-slave' }
    options {
        buildDiscarder logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '', numToKeepStr: '10')
        timeout(time: 90, unit: 'MINUTES')
    }
    environment {
        OPENSHIFT_INTERNAL_REGISTRY = 'image-registry.openshift-image-registry.svc:5000'

        // Use buildah container engine in this pipeline
        CONTAINER_ENGINE = 'podman'
    }
    stages {
        stage('Initialize') {
            steps {
                script {
                    sh ' git config --global user.email "jenkins@kie.com" '
                    sh ' git config --global user.name "kie user"'
                    githubscm.checkoutIfExists('rhpam-kogito-operator', changeAuthor, changeBranch, 'kiegroup', changeTarget, true, ['token' : 'GITHUB_TOKEN', 'usernamePassword' : 'user-kie-ci10'])
                    sh "set +x && oc login --token=\$(oc whoami -t) --server=${OPENSHIFT_API} --insecure-skip-tls-verify"
                    sh """
                        usermod --add-subuids 10000-75535 \$(whoami)
                        usermod --add-subgids 10000-75535 \$(whoami)
                    """
                }
            }
        }
        stage('Test Kogito Operator') {
            steps {
                sh 'make test'
            }
        }
        stage('Build Kogito Operator') {
            steps {
                sh "make BUILDER=${CONTAINER_ENGINE}"
            }
        }
        stage('Push Operator Image to Openshift Registry') {
            steps {
                sh """
                    set +x && ${CONTAINER_ENGINE} login -u jenkins -p \$(oc whoami -t) --tls-verify=false ${OPENSHIFT_REGISTRY}
                    cd version/ && TAG_OPERATOR=\$(grep -m 1 'Version =' version.go) && TAG_OPERATOR=\$(echo \${TAG_OPERATOR#*=} | tr -d '"')
                    ${CONTAINER_ENGINE} tag registry.stage.redhat.io/rhpam-7/rhpam-kogito-rhel8-operator:\${TAG_OPERATOR} ${OPENSHIFT_REGISTRY}/openshift/rhpam-kogito-operator:pr-\$(echo \${GIT_COMMIT} | cut -c1-7)
                    ${CONTAINER_ENGINE} push --tls-verify=false ${OPENSHIFT_REGISTRY}/openshift/rhpam-kogito-operator:pr-\$(echo \${GIT_COMMIT} | cut -c1-7)
                """
            }
        }

        stage('Run BDD tests') {
            options {
                lock("BDD tests ${OPENSHIFT_API}")
            }
            stages {
                stage("Build examples' images for testing") {
                    steps {
                        // Do not build native images for the PR checks
                        sh "make build-smoke-examples-images tags='@rhpam' concurrent=3 ${getBDDParameters('never', false)}"
                    }
                    post {
                        always {
                            archiveArtifacts artifacts: 'test/logs/**/*.log', allowEmptyArchive: true
                            junit testResults: 'test/logs/junit.xml', allowEmptyResults: true
                        }
                    }
                }
                stage('Running smoke tests') {
                    steps {
                        // Run just smoke tests to verify basic operator functionality
                        sh "make run-smoke-tests tags='@rhpam' concurrent=5 ${getBDDParameters('always', true)}"
                    }
                    post {
                        always {
                            archiveArtifacts artifacts: 'test/logs/**/*.log', allowEmptyArchive: true
                            junit testResults: 'test/logs/junit.xml', allowEmptyResults: true
                        }
                    }
                }
            }
        }
    }
    post {
        always {
            cleanWs()
        }
    }
}

String getBDDParameters(String image_cache_mode, boolean runtime_app_registry_internal=false) {
    testParamsMap = [:]

    testParamsMap['load_default_config'] = true
    testParamsMap['ci'] = 'jenkins'
    testParamsMap['load_factor'] = 3
    testParamsMap['disable_maven_native_build_container'] = true

    testParamsMap['operator_image'] = "${OPENSHIFT_REGISTRY}/openshift/rhpam-kogito-operator"
    testParamsMap['operator_tag'] = "pr-\$(echo \${GIT_COMMIT} | cut -c1-7)"

    // Product operator doesn't have CLI
    testParamsMap['cr_deployment_only'] = true

    if (env.MAVEN_MIRROR_REPOSITORY) {
        testParamsMap['maven_mirror'] = env.MAVEN_MIRROR_REPOSITORY
        testParamsMap['maven_ignore_self_signed_certificate'] = true
    }

    // runtime_application_image are built in this pipeline so we can just use Openshift registry for them
    testParamsMap['image_cache_mode'] = image_cache_mode
    testParamsMap['runtime_application_image_registry'] = runtime_app_registry_internal ? env.OPENSHIFT_INTERNAL_REGISTRY : env.OPENSHIFT_REGISTRY
    testParamsMap['runtime_application_image_namespace'] = 'openshift'
    testParamsMap['runtime_application_image_version'] = "pr-\$(echo \${GIT_COMMIT} | cut -c1-7)"

    // Using upstream images as a workaround until there are nightly product images available
    testParamsMap['build_s2i_image_tag'] = "quay.io/kiegroup/kogito-builder-nightly:latest"
    testParamsMap['build_runtime_image_tag'] = "quay.io/kiegroup/kogito-runtime-jvm-nightly:latest"

    testParamsMap['container_engine'] = env.CONTAINER_ENGINE

    String testParams = testParamsMap.collect { entry -> "${entry.getKey()}=\"${entry.getValue()}\"" }.join(' ')
    echo "BDD parameters = ${testParams}"
    return testParams
}
