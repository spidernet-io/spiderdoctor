# Introduction

## copy

1. copy repo

   replace all 'spiderdoctor' to 'YourRepoName'

   replace all 'spidernet-io' and 'spidernet.io' to 'YourOrigin'

2. grep "====modify====" * -RHn --colour  and modify all of them

3. update api/v1/openapi.yaml and `make update_openapi_sdk`

4. redefine CRD in pkg/k8s/v1
    rename directory name 'pkg/k8s/v1/spiderdoctor.spidernet.io' 
    replace all 'mybook' to 'YourCRDName'
    and `make update_crd_sdk`, and code pkg/mybookManager

5. update charts/ , and images/ , and CODEOWNERS

6. go mod tidy , go mod vendor , go vet ./... , double check all is ok

7. create an empty branch 'github_pages' and mkdir 'docs'

8. github seetings:

   spidernet.io  -> settings -> secrets -> actions -> grant secret to repo

   spidernet.io  -> settings -> general -> feature -> issue

   repo -> packages -> package settings -> Change package visibility

   create 'github_pages' branch, and repo -> settings -> pages -> add branch 'github_pages', directory 'docs'

   repo -> settings -> branch -> add protection rules for 'main' and 'github_pages'

9. enable third app

   codefactor: https://www.codefactor.io/dashboard

   sonarCloud: https://sonarcloud.io/projects/create

10. create badge for github/workflows/auto-ci.yaml, github/workflows/badge.yaml

11. build base image , 
    update BASE_IMAGE in images/agent/Dockerfile and images/controller/Dockerfile
    run test


## local develop

1. `make build_local_image`

2. `make e2e_init`

3. `make e2e_run`

4. check proscope, browser vists http://NodeIP:4040

5. apply cr

        cat <<EOF > mybook.yaml
        apiVersion: spiderdoctor.spidernet.io/v1
        kind: Mybook
        metadata:
          name: test
        spec:
          ipVersion: 4
          subnet: "1.0.0.0/8"
        EOF
        kubectl apply -f mybook.yaml

## chart develop

helm repo add rock https://spidernet-io.github.io/spiderdoctor/

