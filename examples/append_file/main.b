// Программа дозаписи в файл 

#import "stdlib/files.b"

void main(){
	int f;
	f = open_file("test.txt", "append");
	write_file(f, "Hello world!\n");
	close_file(f);
};

main();
