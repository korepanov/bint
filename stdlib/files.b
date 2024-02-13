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
	}else if (mode == "append"){
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
 reads size bytes from file with descriptor_number

 returns read info as a string 
*/
string read_file(int descriptor_number, int size){
	int bytes; 
	string s; 

	bytes = $read_f(descriptor_number, size, s);
	return s;
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
		int bytes; 
		int size;
		size = 1;
		bytes = $read_f(descriptor_number, size, s); 
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
