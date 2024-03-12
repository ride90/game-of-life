/* Requires the Docker Pipeline plugin */
pipeline {
    agent { docker { image 'golang:1.16' } }
    stages {
        stage('build') {
            steps {
                sh 'go version'
            }
        }
    }
}