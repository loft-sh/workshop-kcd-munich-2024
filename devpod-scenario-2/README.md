# Data Analysts

## Scenario

We have an existing development environment in our kubernetes cluster.
A new team of data analysts joins and need access to some services on the cluster.

## Setup

Existing virtual clusters should already have the environment set up, run `kubectl get all -n devpod-demo` to confirm.

Example Output:

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

Otherwise, run `kubectl apply devpod-scenario-2/deploy`.

## Next Steps

- node selector
- resources
