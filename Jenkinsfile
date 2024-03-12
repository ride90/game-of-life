/* Requires the Docker Pipeline plugin */
pipeline {
    agent {
        docker {
            image 'golang:1.16'
            label 'docker'
        }
    }
    stages {
        stage('build') {
            steps {
                sh 'go version'
            }
        }
    }
}