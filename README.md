##����beekeeper

һ��GO����ʵ�ֵ�������fcgi��ܡ�

##GO�汾Ҫ��

go1.7�����ϡ�

##��Ҫģ��

* conf
	*  ͳһ�����ü���ģ�飬���Զ�ȡ�����в�����xml�����ļ�������Ϊ����
	
* grace
	*  �ṩ���жϷ����ܣ������������ڼ䲻��Ͽ��������ӣ�ֱ����Щ���ӳ�ʱ��
	
* log
	*  ���ڱ�׼��log��ʵ�ֵ�rotate��־ģ�顣

* mon
	*  ���ڱ�׼��expvar��ʵ�ֵļ�ؽӿ�ģ�飬���Խ����������ڼ��һЩͳ�Ʊ���ͨ��httpЭ����ͳһ�ĸ�ʽ��¶���ⲿ��

* router
	*  һ���򵥵�����·��ģ�飬����ע���·����Ϣ����ͬ·���ͷ���������·�ɵ���Ӧ�ĺ�������

##��δ����Լ��Ĵ���

��ȷ����ȷ��װGO������ǰ���£�����beekeeper/create.bat��create.sh���ű�����$GOPATH/srcĿ¼������Ŀ¼��

�������ɵ�Ŀ¼��ʹ��go build������롣

##����web������

��nginx�������������ת������

```javascript
location /your_program_name/ {
	fastcgi_pass x.x.x.x:yyyy;
	fastcgi_index index.cgi;
	fastcgi_param SCRIPT_FILENAME fcgi$fastcgi_script_name;
	include fastcgi_params;
}
```

ע��ת����ip��ַ�Ͷ˿�Ӧ���������ļ��е�ip��ַ�Ͷ˿ڱ���һ�¡�