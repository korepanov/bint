float sum(float a, float b){
	float res; 
	res = (a + b);
	return res; 
};

float sub(){
	float res;
	float a; 
	float b;
	a = sum(1, 2);
	b = sum(3, 4);
	 
	res = (a - b);
	return res; 
};

float mul(){
	float a; 
	float b; 
	a = sub();
	b = 5; 
	
	res = (a * b);
	return res; 

};

float res; 
res = mul();
string buf;
buf = str(res);
print(buf);  

