pipeline {
    agent any

    environment {
        DOCKER_REGISTRY = 'docker.io/denisgjonmarkaj'
        DOCKER_CREDENTIALS = credentials('dockerhub-credentials')
        BUILD_VERSION = "${BUILD_NUMBER}"
        GOROOT = '/usr/local/go'
        PATH = "${env.GOROOT}/bin:${env.PATH}"
        NODE_OPTIONS = "--openssl-legacy-provider"  // Aggiunto per risolvere l'errore OpenSSL
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
                                export NODE_OPTIONS=--openssl-legacy-provider
                                npm ci
                                npm audit fix --force || true
                            '''
                        }
                    }
                }
                
                stage('Backend') {
                    steps {
                        dir('backend') {
                            sh 'go mod download || true'
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
                                export NODE_OPTIONS=--openssl-legacy-provider
                                export CI=true
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
                                if [ -f main.go ]; then
                                    CGO_ENABLED=0 GOOS=linux go build -o main .
                                    docker build -t ${DOCKER_REGISTRY}/meal-planner-backend:${BUILD_VERSION} .
                                else
                                    echo "No main.go found in backend directory"
                                    ls -la
                                    exit 1
                                fi
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
