pipeline {
    agent any

    stages {
            stage ('Package') {
                steps {
                    sh './gradlew build -x test'
                }
            }
    }
}
