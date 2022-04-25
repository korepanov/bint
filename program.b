#import "stdlib/strings.b"

void hello_user(string s){
	s = (("Hello, " + s) + "!");
	println(s);
};
void main(){
	string s;
	input(s);
	hello_user(s);
};

main();
