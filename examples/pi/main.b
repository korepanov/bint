// Программа вычисление числа Пи методом Мадхава

float calc(int n){
	float t;
	float pi;
	t = 1;
	pi = 1;

	for (int y; y = 3; y < (n*2); y = (y + 2)){
		t = ((-3.0)*t);
		pi = (pi + ((1.0 / t) / float(y)));
	};

	return ((12.0^0.5)*pi);
};
void main(){
	print("Программа вычисления числа Пи методом Мадхава\n");
	float pi;
	pi = calc(10);
	print("Число Пи: ");
	print((str(pi) + "\n"));
};
main();
