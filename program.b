#import "stdlib/strings.b"

void hello_user(string s){
	s = (("Привет, " + s) + "!");
	println(s);
};
void main(){
	string s;
	println("Ваше имя:");
	input(s);
	hello_user(s);
};

main();
