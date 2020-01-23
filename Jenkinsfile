pipeline {
    agent { 
        docker { 
            image 'epel7mgo' 
            label 'slave01'
            args '-u jenkins:jenkins'
        }
    }
    options { checkoutToSubdirectory('argo-web-api') }
    stages {
         
        stage('Build') {
            
            steps {
                echo 'Build...'
                sh '''export GIT_COMMIT=$(git rev-list -1 HEAD)
                export BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
                mkdir -p ${WORKSPACE}/go/src/github.com/ARGOeu
                export GOPATH=${WORKSPACE}/go
                ln -sf ${WORKSPACE}/argo-web-api ${WORKSPACE}/go/src/github.com/ARGOeu/argo-web-api
                rm -rf ${WORKSPACE}/go/src/github.com/ARGOeu/argo-web-api/argo-web-api
                cd ${WORKSPACE}/go/src/github.com/ARGOeu/argo-web-api && go build
                '''
            }
        }
        stage('Test') {
            steps {
                echo 'Test & Coverage...'
                sh '''mkdir -p ${WORKSPACE}/go/src/github.com/ARGOeu
                export GOPATH=${WORKSPACE}/go
                ln -sf ${WORKSPACE}/argo-web-api ${WORKSPACE}/go/src/github.com/ARGOeu/argo-web-api
                sudo /etc/init.d/mongod restart
                cd ${WORKSPACE}/go/src/github.com/ARGOeu/argo-web-api && gocov test -p 1 $(go list ./... | grep -v /vendor/) | gocov-xml > ${WORKSPACE}/coverage.xml
                cd ${WORKSPACE}/go/src/github.com/ARGOeu/argo-web-api && go test -p 1 $(go list ./... | grep -v /vendor/) -v=1 | go-junit-report > ${WORKSPACE}/junit.xml
                '''
                junit '**/junit.xml'
                cobertura coberturaReportFile: '**/coverage.xml'

            }
        }
        stage('Package') {
            steps {
                echo 'Building Rpm...'
                sh '''cd ${WORKSPACE}/argo-web-api && make sources
                cp ${WORKSPACE}/argo-web-api/argo-web-api*.tar.gz /home/jenkins/rpmbuild/SOURCES/
                cd ${WORKSPACE}/argo-web-api && export GIT_COMMIT=${GIT_COMMIT:-`git log -1 --format="%H"`}
                cd ${WORKSPACE}/argo-web-api && export GIT_COMMIT_HASH=`echo ${GIT_COMMIT} | cut -c1-7`
                cd ${WORKSPACE}/argo-web-api && export _GIT_COMMIT_DATE=`git show -s --format=%ci ${GIT_COMMIT_HASH}`
                cd ${WORKSPACE}/argo-web-api && export GIT_COMMIT_DATE=`date -d "${_GIT_COMMIT_DATE}" "+%Y%m%d%H%M%S"`
                sed -i 's/^Release.*/Release: %(echo $GIT_COMMIT_DATE).%(echo $GIT_COMMIT_HASH)%{?dist}/' ${WORKSPACE}/argo-web-api/argo-web-api.spec
                cd /home/jenkins/rpmbuild/SOURCES && tar -xzvf argo-web-api*.tar.gz
                cp ${WORKSPACE}/argo-web-api/argo-web-api.spec /home/jenkins/rpmbuild/SPECS/
                rpmbuild -bb /home/jenkins/rpmbuild/SPECS/*.spec
                rm -f ${WORKSPACE}/*.rpm
                cp /home/jenkins/rpmbuild/RPMS/**/*.rpm ${WORKSPACE}/
                '''
                archiveArtifacts artifacts: '**/*.rpm', fingerprint: true
            }
        }
       
        stage('Package for Devel') {
            when {
                branch 'devel'
            }
            steps {
                echo 'Building Rpm for Devel...'
                sh '''cd ${WORKSPACE}/argo-web-api && make sources
                cp ${WORKSPACE}/argo-web-api/argo-web-api*.tar.gz /home/jenkins/rpmbuild/SOURCES/
                cd ${WORKSPACE}/argo-web-api && export GIT_COMMIT=${GIT_COMMIT:-`git log -1 --format="%H"`}
                cd ${WORKSPACE}/argo-web-api && export GIT_COMMIT_HASH=`echo ${GIT_COMMIT} | cut -c1-7`
                cd ${WORKSPACE}/argo-web-api && export _GIT_COMMIT_DATE=`git show -s --format=%ci ${GIT_COMMIT_HASH}`
                cd ${WORKSPACE}/argo-web-api && export GIT_COMMIT_DATE=`date -d "${_GIT_COMMIT_DATE}" "+%Y%m%d%H%M%S"`
                sed -i 's/^Release.*/Release: %(echo $GIT_COMMIT_DATE).%(echo $GIT_COMMIT_HASH)%{?dist}/' ${WORKSPACE}/argo-web-api/argo-web-api.spec
                cd /home/jenkins/rpmbuild/SOURCES && tar -xzvf argo-web-api*.tar.gz
                cp ${WORKSPACE}/argo-web-api/argo-web-api.spec /home/jenkins/rpmbuild/SPECS/
                rpmbuild -bb /home/jenkins/rpmbuild/SPECS/*.spec
                rm -f ${WORKSPACE}/*.rpm
                cp /home/jenkins/rpmbuild/RPMS/**/*.rpm ${WORKSPACE}/
                '''
                archiveArtifacts artifacts: '**/*.rpm', fingerprint: true
                echo 'Uploading rpm for devel...'
                withCredentials(bindings: [sshUserPrivateKey(credentialsId: 'jenkins-repo', usernameVariable: 'REPOUSER', \
                                                             keyFileVariable: 'REPOKEY')]) {
                  sh '''
                  scp -i ${REPOKEY} -o StrictHostKeyChecking=no ${WORKSPACE}/*.rpm ${REPOUSER}@rpm-repo.argo.grnet.gr:/repos/ARGO/devel/centos7/
                  '''
                }
            }
        }
        stage('Package for Production') {
            when {
                branch 'master'
            }
            steps {
                echo 'Building Rpm for Production...'
                sh '''cd ${WORKSPACE}/argo-web-api && make sources
                cd /home/jenkins/rpmbuild/SOURCES && tar -xzvf argo-web-api*.tar.gz
                cp ${WORKSPACE}/argo-web-api/argo-web-api.spec /home/jenkins/rpmbuild/SPECS/
                rpmbuild -bb /home/jenkins/rpmbuild/SPECS/*.spec
                rm -f ${WORKSPACE}/*.rpm
                cp /home/jenkins/rpmbuild/RPMS/**/*.rpm ${WORKSPACE}/
                '''
                archiveArtifacts artifacts: '**/*.rpm', fingerprint: true
                echo 'Uploading rpm for devel...'
                withCredentials(bindings: [sshUserPrivateKey(credentialsId: 'jenkins-repo', usernameVariable: 'REPOUSER', \
                                                             keyFileVariable: 'REPOKEY')]) {
                  sh '''
                  scp -i ${REPOKEY} -o StrictHostKeyChecking=no ${WORKSPACE}/*.rpm ${REPOUSER}@rpm-repo.argo.grnet.gr:/repos/ARGO/prod/centos7/
                  '''
                }
            }
        }
        
    }
}
