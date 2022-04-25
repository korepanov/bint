#import "stdlib/strings.b"

void hello_user(string s){
	s = (("Hello, " + s) + "!");
	println(s);
};
void main(){
	hello_user("Slava");
};

main();
