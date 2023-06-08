float find_descr(float a, float b, float c){
	return ((b^2.0)-(4.0*(a*c)));
};


void main(){
	print("Решение квадратного уравнения x^2 + 2x - 3\n");
	float descr;
	descr = find_descr(1.0, "lalala", -3.0); 
	print(str(descr));
	print("\n");

};

main();
