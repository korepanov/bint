float find_x1(float descr, float a, float b, float c){
	return ((((-1)*b)-(descr^0.5)) / (2*a)); 
};

float find_descr(float a, float b, float c){
	return ((b^2.0)-(4.0*(a*c)));
};


void main(){
	print("Решение квадратного уравнения x^2 + 2x - 3\n");
	float x1;
	x1 = find_x1(find_descr(1.0, 2.0, -3.0), 1.0, 2.0, -3.0); 
	print(str(x1));
	print("\n");

};

main();
