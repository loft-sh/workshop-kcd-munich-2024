# Data Analysts

## Scenario

We have an existing development environment in our kubernetes cluster.
A new team of data analysts joins and need access to some services on the cluster.

## Setup

1. Create a local cluster (kind, minikube, docker desktop, orbstack, ...)
2. Run `kubectl apply -f devpod-scenario-2/deploy` to initialize the cluster. It'll deploy a mysql instance and a phpadmin dashboard


Output of `kubectl get all -n devpod-demo` to confirm:
```
NAME                              READY   STATUS    RESTARTS   AGE
pod/mysql-5866b56598-qvvxm        1/1     Running   0          140m
pod/phpmyadmin-67dbb8db8d-vncrr   1/1     Running   0          141m

NAME                 TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)           AGE
service/mysql        ClusterIP      None           <none>        3306/TCP          140m
service/phpmyadmin   LoadBalancer   10.103.142.0   <pending>     12080:31689/TCP   141m

NAME                         READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/mysql        1/1     1            1           140m
deployment.apps/phpmyadmin   1/1     1            1           141m

NAME                                    DESIRED   CURRENT   READY   AGE
replicaset.apps/mysql-5866b56598        1         1         1       140m
replicaset.apps/phpmyadmin-67dbb8db8d   1         1         1       141m
```
