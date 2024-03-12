/* Requires the Docker Pipeline plugin */
pipeline {
    agent {
        label 'docker'
    }
    stages {
        stage('build') {
            agent {
                docker {
                    image 'golang:1.16'
                    label 'docker'
                }
            }
            steps {
                sh 'go version'
            }
        }
    }
}