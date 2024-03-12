/* Requires the Docker Pipeline plugin */
pipeline {
    agent { docker { image 'golang:1.22.1-alpine3.19' } }
    stages {
        stage('build') {
            steps {
                sh 'go version'
            }
        }
    }
}