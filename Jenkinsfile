pipeline {
    agent { 
        docker { 
            image 'argo.registry:5000/rocky9-mongo6.0-go1.22' 
            args '-u jenkins:jenkins'
        }
    }
    options {
        checkoutToSubdirectory('argo-web-api')
        newContainerPerStage()
    }
    environment {
        PROJECT_DIR='argo-web-api'
        GH_USER = 'newgrnetci'
        GH_EMAIL = '<argo@grnet.gr>'
        GIT_COMMIT=sh(script: "cd ${WORKSPACE}/$PROJECT_DIR && git log -1 --format=\"%H\"",returnStdout: true).trim()
        GIT_COMMIT_HASH=sh(script: "cd ${WORKSPACE}/$PROJECT_DIR && git log -1 --format=\"%H\" | cut -c1-7",returnStdout: true).trim()
        GIT_COMMIT_DATE=sh(script: "date -d \"\$(cd ${WORKSPACE}/$PROJECT_DIR && git show -s --format=%ci ${GIT_COMMIT_HASH})\" \"+%Y%m%d%H%M%S\"",returnStdout: true).trim()
   }
    stages {
        stage('Build') {
            steps {
                echo 'Build...'
                sh """
                cd ${WORKSPACE}/${PROJECT_DIR}
                export CGO_CFLAGS"=-O2 -fstack-protector --param=ssp-buffer-size=4 -D_FORTIFY_SOURCE=2"
                go build -buildmode=pie -ldflags "-s -w -linkmode=external -extldflags '-z relro -z now'"
                """

                archiveArtifacts artifacts: '**/argo-web-api'
            }
        }
        stage('Test') {
            steps {
                echo 'Test & Coverage...'
                sh """
                sudo mongod --fork --logpath /var/log/mongodb.log 
                cd ${WORKSPACE}/${PROJECT_DIR}
                go mod tidy
                gocov test -p 1 ./... | gocov-xml > ${WORKSPACE}/coverage.xml
                go test -p 1 ./... -v=1 | go-junit-report > ${WORKSPACE}/junit.xml
                """
                junit '**/junit.xml'
                cobertura coberturaReportFile: '**/coverage.xml'

            }
        }
    }
    post{
        always {
            cleanWs()
        }
        success {
            script{
                if ( env.BRANCH_NAME == 'devel' ) {
                    build job: '/ARGO/argodoc/devel', propagate: false
                } else if ( env.BRANCH_NAME == 'master' ) {
                    build job: '/ARGO/argodoc/master', propagate: false
                }
                if ( env.BRANCH_NAME == 'master' || env.BRANCH_NAME == 'devel' ) {
                    slackSend( message: ":rocket: New version for <$BUILD_URL|$PROJECT_DIR>:$BRANCH_NAME Job: $JOB_NAME !")
                }
            }
        }
        failure {
            script{
                if ( env.BRANCH_NAME == 'master' || env.BRANCH_NAME == 'devel' ) {
                    slackSend( message: ":rain_cloud: Build Failed for <$BUILD_URL|$PROJECT_DIR>:$BRANCH_NAME Job: $JOB_NAME")
                }
            }
        }
    }
}