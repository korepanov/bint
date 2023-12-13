#import "stdlib/files.b"

void main(){
	int f;
	f = open_file("home/slava/a1.txt", "read");
	string s;
	int bytes;
	bytes = read_file(f, 10, s);
	close_file(f);
	print((s + "\n")); 
};

main();
