 ### helm install prom prometheus-community/kube-prometheus-stack

 #### And for docker deskstop remove volumemounts 

 ### kubectl patch ds prom-prometheus-node-exporter --type "json" -p ' [ {"op": "remove", "path" : "/spec/template/spec/containers/0/volumeMounts/2/mountPropagation"}]'

 ### helm install mongo-exporter prometheus-community/prometheus-mongodb-exporter -f values.yaml