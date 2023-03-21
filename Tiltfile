load("ext://ko", "ko_build")

ko_build("hb", ".")
k8s_yaml(helm('deploy', set=['image.repository=hb']))
k8s_resource('hb', labels=["services"])