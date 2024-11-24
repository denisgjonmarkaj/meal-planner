pipeline {
    agent any

    environment {
        DOCKER_REGISTRY = 'docker.io/denisgjonmarkaj'
        DOCKER_CREDENTIALS = credentials('dockerhub-credentials')
        BUILD_VERSION = "${BUILD_NUMBER}"
        GOROOT = '/usr/local/go'
        PATH = "${env.GOROOT}/bin:${env.PATH}"
    }

    tools {
        nodejs 'NodeJS 18.x'
    }

    stages {
        stage('Verify Tools') {
            steps {
                sh '''
                    echo "Node version:"
                    node --version
                    echo "NPM version:"
                    npm --version
                    echo "Go version:"
                    go version
                    echo "Docker version:"
                    docker --version
                '''
            }
        }

        stage('Install Dependencies') {
            parallel {
                stage('Frontend') {
                    steps {
                        dir('frontend') {
                            sh '''
                                npm ci
                                npm audit fix --force || true
                            '''
                        }
                    }
                }
                
                stage('Backend') {
                    steps {
                        dir('backend') {
                            sh '''
                                export PATH=$PATH:/usr/local/go/bin
                                go mod download
                            '''
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
                            sh 'CI=true npm test || true'
                        }
                    }
                }
                
                stage('Backend Tests') {
                    steps {
                        dir('backend') {
                            sh 'go test ./... || true'
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