pipeline {
    agent any

    stages {
        stage("Push to server") {
            steps {
                git 'gitea@development.alania.lan:Server/recipientverification.rwth.git'

                withCredentials([
                    sshUserPrivateKey(credentialsId: "jenkins-deploy-key", keyFileVariable: 'keyfile'),
                    file(credentialsId: 'recipient-verification-environment', variable: 'environmentFile')
                ]) {
                    sh "ssh -i ${keyfile} jenkinsdeploy@${params.deployHost} 'rm -rf /var/www/recipientVerification/*'"
                    sh "scp -r -i ${keyfile} $WORKSPACE/* jenkinsdeploy@${params.deployHost}:/var/www/recipientVerification"
                    sh "scp -i ${keyfile} ${environmentFile} jenkinsdeploy@${params.deployHost}:/var/www/recipientVerification/.env"
                }
            }
        }

        stage("Run docker-compose up on remote") {
            steps {
                withCredentials([sshUserPrivateKey(credentialsId: "jenkins-deploy-key", keyFileVariable: 'keyfile')]) {
                    sh "ssh -i ${keyfile} jenkinsdeploy@${params.deployHost} 'cd /var/www/recipientVerification/ && docker-compose up --build -d'"
                }
            }
        }
    }
}