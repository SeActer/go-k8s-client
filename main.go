package main

import (
	"flag"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset

func main() {
	k8sconfig := flag.String("k8sconfig", "config", "kubernetes config file path")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *k8sconfig)
	if err != nil {
		log.Println(err)
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("connect k8s success")
	}
	//获取namespace下的所有pod

	namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for _, nsList := range namespaces.Items {
			fmt.Println("##################")
			fmt.Printf("NsName: %s \n", nsList.Name)
			//列出svc后的pod
			services, err := clientset.CoreV1().Services(nsList.Name).List(metav1.ListOptions{})
			if err != nil {
				log.Println(err.Error)
			}
			for _, service := range services.Items {
				if service.Name == "kubernetes" {
					continue
				}
				fmt.Printf("serviceName: %s \n", service.Name)
				fmt.Printf("serviceIp: %s \n", service.Spec.ClusterIP)
				var selectorstring string
				for svcselectorkey, svcselectorvalue := range service.Spec.Selector {
					fmt.Printf("svcselectorkey: %s \n", svcselectorkey)
					fmt.Printf("svcselectorvalue: %s \n", svcselectorvalue)
					selectorstring = svcselectorkey + "=" + svcselectorvalue + "," + selectorstring
				}
				//列出这个deploymentx下的pod,使用标签选择
				fmt.Printf("selectorstring: %s \n", selectorstring[0:len(selectorstring)-1])
				pods, err := clientset.CoreV1().Pods(nsList.Name).List(metav1.ListOptions{LabelSelector: selectorstring[0 : len(selectorstring)-1]})
				if err != nil {
					log.Println(err.Error())
				}
				for _, pod := range pods.Items {
					fmt.Printf("podName: %s \n", pod.Name)
				}
			}

			// deployments, err := clientset.AppsV1().Deployments(nsList.Name).List(metav1.ListOptions{})
			// if err != nil {
			// 	log.Println(err.Error())
			// }
			// for _, deployment := range deployments.Items {
			// 	fmt.Printf("DeploymentName: %s \n", deployment.Name)
			// 	fmt.Printf("DeploymentCreationTimestamp: %s \n", deployment.CreationTimestamp)
			// 	//labelall := list.New()
			// 	var all string
			// 	for labelkey, labelvalue := range deployment.Spec.Selector.MatchLabels {
			// 		fmt.Printf("labelkey: %s \n", labelkey)
			// 		fmt.Printf("labelvalue: %s \n", labelvalue)
			// 		//fmt.Printf(labelkey + "=" + labelvalue)
			// 		all = labelkey + "=" + labelvalue + "," + all
			// 		//将key.value，拼接出来，删掉最后一个,就是需要的字符串
			// 		// labelall.PushBack(labelkey + "=" + labelvalue)
			// 		// for i := labelall.Front(); i != nil; i = i.Next() {
			// 		// 	fmt.Println(i.Value)
			// 		// }
			// 	}
			// 	// for i := labelall.Front(); i != nil; i = i.Next() {
			// 	// 	fmt.Println(i.Value)
			// 	// }
			// 	//列出这个deploymentx下的pod,使用标签选择
			// 	pods, err := clientset.CoreV1().Pods(nsList.Name).List(metav1.ListOptions{LabelSelector: all[0 : len(all)-1]})
			// 	if err != nil {
			// 		log.Println(err.Error())
			// 	}
			// 	for _, pod := range pods.Items {
			// 		fmt.Printf("podName: %s \n", pod.Name)
			// 	}
			// }

			//listpod
			// pods, err := clientset.CoreV1().Pods(nsList.Name).List(metav1.ListOptions{})
			// if err != nil {
			// 	log.Println(err.Error())
			// }
			// for _, pod := range pods.Items {
			// 	fmt.Printf("podName: %s \n", pod.Name)
			// 	fmt.Printf("podCreationTimestamp: %s \n", pod.CreationTimestamp)
			// 	fmt.Printf("podLabels: %s \n", pod.Labels)
			// 	fmt.Printf("podNamespace: %s \n", pod.Namespace)
			// 	fmt.Printf("uuid: %s \n", pod.UID)
			// 	fmt.Printf("HostIP: %s \n", pod.Status.HostIP)
			// 	fmt.Printf("PodIP: %s \n", pod.Status.PodIP)
			// 	fmt.Printf("StartTime: %s \n", pod.Status.StartTime)
			// 	fmt.Printf("Phase: %s \n", pod.Status.Phase)
			// 	for _, containers := range pod.Status.ContainerStatuses {
			// 		fmt.Printf("Image: %s \n", containers.Image)
			// 		fmt.Printf("RestartCount: %d \n", containers.RestartCount)
			// 		fmt.Printf("containerID: %s \n", containers.ContainerID)
			// 	}
			// }
		}
	}

	//
	////获取NODE
	//fmt.Println("##################")
	//nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	//fmt.Println(nodes.Items[0].Name)
	//fmt.Println(nodes.Items[0].CreationTimestamp)    //加入集群时间
	//fmt.Println(nodes.Items[0].Status.NodeInfo)
	//fmt.Println(nodes.Items[0].Status.Conditions[len(nodes.Items[0].Status.Conditions)-1].Type)
	//fmt.Println(nodes.Items[0].Status.Allocatable.Memory().String())

}
