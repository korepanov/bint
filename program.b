#import "stdlib/files.b"

void main(){
	int f;
	
	f = open_file("test.txt", "write");
	write_file(f, "Привет!\n");
	close_file(f); 
};

main();
