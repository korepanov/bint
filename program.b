#import "stdlib/files.b"

void main(){
	int f;
	f = open_file("/home/slava/a1.txt", "read");
	string s;
	s = read_file(f, 10);
	close_file(f);
	print((s + "\n")); 
};

main();
