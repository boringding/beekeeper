##����beekeeper

һ��Go����ʵ�ֵ�������FCGI��ܡ�

##Go�汾Ҫ��

go1.7�����ϡ�

##��Ҫģ��

* conf
	*  ͳһ�����ü���ģ�飬���Զ�ȡ�����в�����XML�����ļ�������Ϊ����
	
* grace
	*  �ṩ���жϷ����ܣ������������ڼ䲻��Ͽ��������ӣ�ֱ����Щ���ӳ�ʱ��
	
* log
	*  ���ڱ�׼��log��ʵ�ֵ�rotate��־ģ�顣

* mon
	*  ���ڱ�׼��expvar��ʵ�ֵļ�ؽӿ�ģ�飬���Խ����������ڼ��һЩͳ�Ʊ����������ڴ�ʹ����������������Լ������ʱ�ȣ�ͨ��HTTPЭ����ͳһ��JSON��ʽ��¶���ⲿ��

* router
	*  һ���򵥵�����·��ģ�飬����ע���·����Ϣ����ͬ·���ͷ���������·�ɵ���Ӧ�ĺ�������

##��δ����Լ��Ĵ���

��Go�����Ѿ���ȷ��װ�����õ�ǰ���£�

1. ��beekeeper��ܸ��Ƶ�$GOPATH/srcĿ¼�¡�

2. ����beekeeper/create.bat��create.sh���ű�����$GOPATH/srcĿ¼������Ŀ¼��

3. �������ɵ�Ŀ¼��ʹ��go build������롣

##����web������

��nginx�������������ת������

```nginx
location /your_program_name/ {
	fastcgi_pass x.x.x.x:yyyy;
	fastcgi_index index.cgi;
	fastcgi_param SCRIPT_FILENAME fcgi$fastcgi_script_name;
	include fastcgi_params;
}
```

ע��ת����IP��ַ�Ͷ˿�Ӧ���������ļ��е�IP��ַ�Ͷ˿ڱ���һ�¡�