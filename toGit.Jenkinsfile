pipeline {
    agent { label 'jenkins-agent-cloud' }
    stages {
        stage('Mirror to Github') {
            steps {
                checkout([
                    $class: 'GitSCM',
                    branches: [[name: '*/master']],
                    extensions: [
                        [$class: 'CleanBeforeCheckout'],
                    ],
                    userRemoteConfigs: [[
                        url: 'https://gerrit.corp.arista.io/openconfigbeat',
                    ]],
                ])
                sshagent (credentials: ['jenkins-rsa-key']) {
                    sh 'if [ ! -d "~/.ssh" ]; then mkdir ~/.ssh; fi'
                    sh 'if ! grep -q github.com ~/.ssh/known_hosts; then ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts; fi'
                    sh 'git push git@github.com:aristanetworks/openconfigbeat.git HEAD:master'
                }
            }
        }
    }
}
