# geoserver

# building
1. make push

# running
1. make sure `kubefed` and `kubectl` are in `$PATH`
1. create cluster

    ```
    ./create.sh
    ```
    
1. deploy service, ingress and replica set

    ```
    kubectl --context=federation create -f services/geoserver.yaml
    kubectl --context=federation create -f rs/geoserver.yaml
    ```
    
1. open the `test.html` page
