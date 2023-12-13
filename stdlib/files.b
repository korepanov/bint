/* opens file with file_path and mode 
 
 modes:
 mode = read 
 mode = write
 mode = append 

 return file descriptor number  
*/

int open_file(string file_path, string mode){
	int res; 
	int m;

	if (mode == "read"){
		m = 0;
	}else if (mode == "write"){
		m = 1;	
	}else if (mode = "append"){
		m = 2;	
	}else{
		error = "no such mode for opening file";
		if ($toPanic){
			print((error + "\n"));
			exit(1);		
		};
		return res; 
	};

	res = $open_f(file_path, m); 
	return res;
};

/*
 reads size bytes from file with descriptor_number into variable s

 returns number of read bytes 
*/
int read_file(int descriptor_number, int size; string s){
	int bytes; 
	bytes = $read_f(descriptor_number, size, s);
	return bytes;
};

/*
 closes file with descriptor_number 
*/
void close_file(int descriptor_number){
	bool $toPanicOld;
	$toPanicOld = $toPanic;
	
	$close_f(descriptor_number); 
	
	try{
		string s; 
		read_file(descriptor_number, 1, s);	
	};
	if (error == ""){
		error = ("could not close file with file descriptor number " + str(descriptor_number));	
		
		if ($toPanicOld){
			print((error + "\n"));
			exit(1);		
		};	
	};
	error = "";
};
