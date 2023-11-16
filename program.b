float get_descr(float a, float b, float c){
	return (((b^2.0)-((4*a)*c)));
};

void solve(float a, float b, float c, float d){
	if (d < 0){
		print("Нет решений\n");	
	}else{
		float x1;
		float x2;
		
		x1 = ((((-1.0)*b)-(d^0.5))/(2.0*a));
		x2 = ((((-1.0)*b)+(d^0.5))/(2.0*a));

		print((("x1=" + str(x1)) + "\n"));	
		print((("x2=" + str(x2)) + "\n"));
	};
};

void main(){
	print("Решение квадратных уравнений вида ax^2+bx+c=0\n");
	float a;
	float b;
	float c;
	string s;
	print("Введите a\n");
	input(s);
	a = float(s); 
	print("Введите b\n");
	input(s);
	b = float(s); 
	print("Введите c\n");
	input(s);
	c = float(s); 
	float d;
	d = get_descr(a, b, c); 
	solve(a, b, c, d);

};

main();
