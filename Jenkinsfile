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
        stage('Build') {
            steps {
                echo 'Build...'
                sh """
                mkdir -p ${WORKSPACE}/go/src/github.com/ARGOeu
                ln -sf ${WORKSPACE}/${PROJECT_DIR} ${WORKSPACE}/go/src/github.com/ARGOeu/${PROJECT_DIR}
                rm -rf ${WORKSPACE}/go/src/github.com/ARGOeu/${PROJECT_DIR}/${PROJECT_DIR}
                cd ${WORKSPACE}/go/src/github.com/ARGOeu/${PROJECT_DIR}
                go build
                """
            }
        }
        stage('Test') {
            steps {
                echo 'Test & Coverage...'
                sh """
                mkdir /home/jenkins/mongo_data
                mkdir /home/jenkins/mongo_log
                mkdir /home/jenkins/mongo_run
                mongod --dbpath /home/jenkins/mongo_data --logpath /home/jenkins/mongo_log/mongo.log --pidfilepath /home/jenkins/mongo_run/mongo.pid --fork
                cd ${WORKSPACE}/go/src/github.com/ARGOeu/${PROJECT_DIR}
                gocov test -p 1 \$(go list ./... | grep -v /vendor/) | gocov-xml > ${WORKSPACE}/coverage.xml
                go test -p 1 \$(go list ./... | grep -v /vendor/) -v=1 | go-junit-report > ${WORKSPACE}/junit.xml
                """
                junit '**/junit.xml'
                cobertura coberturaReportFile: '**/coverage.xml'

            }
        }
        stage('Package') {
            steps {
                echo 'Building Rpm...'
                withCredentials(bindings: [sshUserPrivateKey(credentialsId: 'jenkins-rpm-repo', usernameVariable: 'REPOUSER', \
                                                             keyFileVariable: 'REPOKEY')]) {
                    sh "/home/jenkins/build-rpm.sh -w ${WORKSPACE} -b ${BRANCH_NAME} -d centos7 -p ${PROJECT_DIR} -s ${REPOKEY}"
                }
                archiveArtifacts artifacts: '**/*.rpm', fingerprint: true
            }
        }
        stage('Deploy to devel') {
            agent { 
                docker { 
                    image 'node:buster' 
                }
            }
            steps {
                echo 'Deploying to development environment'
                withCredentials([sshUserPrivateKey(credentialsId: 'jenkins-rpm-repo', usernameVariable: 'REPOUSER', \
                                                             keyFileVariable: 'REPOKEY'),
                                usernamePassword(credentialsId: 'jenkins-rpm-repo',usernameVariable: 'ARGO_TOKEN', passwordVariable: 'TOKEN_VAL')]) {
                    sh """
                        cd ${WORKSPACE}/${PROJECT}
                        npm install newman
                        cat ./deploy.sh | ssh -i ${REPOKEY} -o StrictHostKeyChecking=no root@snf-13121.ok-kno.grnetcloud.net DEP_PROJECT=argo-web-api DEP_VERSION=1.9.0 DEP_RELEASE=20200402153255.8e1fd9e.el7 /bin/bash
                        ./node_modules/newman/bin/newman.js run ./postman/web-api.postman_collection.json  -k -e ./postman/environment.json --env-var last_commit=8e1fd9e9040f266f90c9fc0be645ae600c3bfd04 --env-var api_key=${TOKEN_VAL}
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
