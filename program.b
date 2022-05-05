#import "stdlib/strings.b"

void hello_user(string s){
	s = (("Hello, " + s) + "!");
	println(s);
};
void main(){
	string s;
	println("Your name:");
	input(s);
	hello_user(s);
};

main();
