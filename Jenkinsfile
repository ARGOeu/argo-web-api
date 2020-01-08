pipeline {
    agent { 
        docker { 
            image 'epel7go' 
            label 'slave02'
            args '-u root:root'
        }
    }
    options { checkoutToSubdirectory('argo-web-api') }
    stages {
        stage('Build') {
            
            steps {
                echo 'Building...'
                sh '''export GIT_COMMIT=$(git rev-list -1 HEAD)
                export BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
                mkdir -p ${WORKSPACE}/go/src/github.com/ARGOeu
                export GOPATH=${WORKSPACE}/go
                ln -sf ${WORKSPACE}/argo-web-api ${WORKSPACE}/go/src/github.com/ARGOeu/argo-web-api
                cd ${WORKSPACE}/go/src/github.com/ARGOeu/argo-web-api && go build
                '''
            }
        }
        stage('Test') {
            steps {
                echo 'Testing ...'
                sh '''mkdir -p ${WORKSPACE}/go/src/github.com/ARGOeu
                export GOPATH=${WORKSPACE}/go
                ln -sf ${WORKSPACE}/argo-web-api ${WORKSPACE}/go/src/github.com/ARGOeu/argo-web-api
                /etc/init.d/mongod restart
                cd ${WORKSPACE}/go/src/github.com/ARGOeu/argo-web-api && go test ./... -p=1
                '''

            }
        }
        stage('Deploy Devel') {
            when {
                branch 'devel'
            }
            steps {
                echo 'Deploying to Devel....'
            }
        }
        stage('Deploy Production') {
            when {
                branch 'master'
            }
            steps {
                echo 'Deploying to Master....'
            }
        }
    }
    post { 
        always { 
            cleanWs()
        }
    }
}
