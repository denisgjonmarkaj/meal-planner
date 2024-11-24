pipeline {
    agent any

    environment {
        DOCKER_REGISTRY = 'docker.io/denisgjonmarkaj'
        DOCKER_CREDENTIALS = credentials('dockerhub-credentials')
        BUILD_VERSION = "${BUILD_NUMBER}"
    }

    tools {
        nodejs 'NodeJS 18.x'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Install Dependencies') {
            parallel {
                stage('Frontend') {
                    steps {
                        dir('frontend') {
                            sh 'npm ci'
                        }
                    }
                }
                
                stage('Backend') {
                    steps {
                        dir('backend') {
                            sh 'go mod download'
                        }
                    }
                }
            }
        }

        stage('Test') {
            parallel {
                stage('Frontend Tests') {
                    steps {
                        dir('frontend') {
                            sh 'npm test -- --watchAll=false'
                        }
                    }
                }
                
                stage('Backend Tests') {
                    steps {
                        dir('backend') {
                            sh 'go test ./...'
                        }
                    }
                }
            }
        }

        stage('Build') {
            parallel {
                stage('Build Frontend') {
                    steps {
                        dir('frontend') {
                            sh '''
                            npm run build
                            docker build -t ${DOCKER_REGISTRY}/meal-planner-frontend:${BUILD_VERSION} .
                            '''
                        }
                    }
                }
                
                stage('Build Backend') {
                    steps {
                        dir('backend') {
                            sh '''
                            CGO_ENABLED=0 GOOS=linux go build -o main .
                            docker build -t ${DOCKER_REGISTRY}/meal-planner-backend:${BUILD_VERSION} .
                            '''
                        }
                    }
                }
            }
        }

        stage('Push Images') {
            steps {
                sh '''
                echo "${DOCKER_CREDENTIALS_PSW}" | docker login -u "${DOCKER_CREDENTIALS_USR}" --password-stdin
                docker push ${DOCKER_REGISTRY}/meal-planner-frontend:${BUILD_VERSION}
                docker push ${DOCKER_REGISTRY}/meal-planner-backend:${BUILD_VERSION}
                '''
            }
        }
    }

    post {
        always {
            sh 'docker logout'
            cleanWs()
        }
        success {
            echo 'Pipeline completed successfully!'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}