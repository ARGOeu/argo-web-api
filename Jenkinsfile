pipeline {
    agent { 
        docker { 
            image 'argo.registry:5000/epel-7-mgo' 
            args '-u jenkins:jenkins'
        }
    }
    options {
        checkoutToSubdirectory('argo-web-api')
        newContainerPerStage()
    }
    environment {
        PROJECT_DIR='argo-web-api'
        GOPATH="${WORKSPACE}/go"
        GIT_COMMIT=sh(script: "cd ${WORKSPACE}/$PROJECT_DIR && git log -1 --format=\"%H\"",returnStdout: true).trim()
        GIT_COMMIT_HASH=sh(script: "cd ${WORKSPACE}/$PROJECT_DIR && git log -1 --format=\"%H\" | cut -c1-7",returnStdout: true).trim()
        GIT_COMMIT_DATE=sh(script: "date -d \"\$(cd ${WORKSPACE}/$PROJECT_DIR && git show -s --format=%ci ${GIT_COMMIT_HASH})\" \"+%Y%m%d%H%M%S\"",returnStdout: true).trim()
   }
    stages {
        stage ('Build messaging docs'){
            environment {
                DOC_SOURCE="argo-messaging"
            }
            steps {
                dir ("${WORKSPACE}/${DOC_SOURCE}") {
                    git branch: "devel",
                        credentialsId: 'jenkins-rpm-repo',
                        url: "git@github.com:ARGOeu/${DOC_SOURCE}.git"
                    sh """
                        cd ${WORKSPACE}/${DOC_SOURCE}/doc/v1
                    """
                    deleteDir()
                }
            }
        }
        stage('Deploy to devel') {
            agent { 
                docker { 
                    image 'argo.registry:5000/epel-7-ams' 
                }
            }
            steps {
                dir ("${WORKSPACE}/grnet-ansible") {
                    git branch: "feature/DEVOPS-138",
                        credentialsId: 'kevangel79',
                        url: "git@github.com:/kevangel79/argo-ansible-deploy.git"
                    sh """
                        cd ${WORKSPACE}/grnet-ansible
                        git submodule init
                        git submodule update
                        cd argo-ansible
                        git fetch origin/feature/DEVOPS-138
                        git checkout feature/DEVOPS-138
                        cd ..
                        pipenv --python 2
                        pipenv run pip install -r argo-ansible/requirements.txt
                        pipenv run ansible-galaxy install -r argo-ansible/requirements.yml
                        echo ">>> Run ansible swagger role"
                        pipenv run ansible-playbook --private-key=${REPO_KEY} -i devel -l testVm13121 argo-ansible/update.yml -u root -vv
                    """
                    deleteDir()
                }
            }
        }
        stage('Run functional tests') {
            agent { 
                docker { 
                    image 'node:buster' 
                }
            }
            steps {
                echo 'Run functional tests at devel environment'
                withCredentials([usernamePassword(credentialsId: 'argo-token',usernameVariable: 'ARGO_TOKEN', passwordVariable: 'TOKEN_VAL')]) {
                    sh """
                        cd ${WORKSPACE}/${PROJECT_DIR}
                        npm install newman
                        ./node_modules/newman/bin/newman.js run ./postman/web-api.postman_collection.json  -k -e ./postman/environment.json --env-var last_commit=9bbfbd575cd79d3dbf246c017e6480bf50318fcc --env-var api_key=${TOKEN_VAL}
                    """
                }
            }
        }
    }
    post{
        success {
            script{
                if ( env.BRANCH_NAME == 'devel' ) {
                    build job: '/ARGO-utils/argo-swagger-docs', propagate: false
                }
            }
        }
    }
}