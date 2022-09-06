#import "stdlib/core.b"

string root_source;

void init(){
	get_root_source(root_source); 
	SET_SOURCE(root_source);
	SET_DEST("benv/dowhile_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

bool is_do(string command){
	stack s;
	string buf;

	buf = "do{";
	s = ops(command, buf);
	s.pop(buf);
	
	return (NOT("end" == buf));
};

string get_cond(string command){
	string op;
	stack s;
	string buf;
	int start_pos;
	int end_pos; 

	op = "while(";
	s = ops(command, op);
	s.pop(buf);
	start_pos = int(buf);
	start_pos = (start_pos + 5);
	end_pos = func_end(command, start_pos);
	start_pos = (start_pos + 1); 
	buf = command[start_pos:end_pos];

	return buf;
};

void main(){
	string buf;
	string cond; 
	int counter;

	init();
	switch_command();

	for (print(""); print(""); NOT(command == "end"); switch_command()){
		if (is_do(command)){
			counter = block_end(); 
			cond = get_cond(get_command(counter));

			buf = "for (print(\"\")";
			send_command(buf); 
			buf = "print(\"\")";
			send_command(buf);
			buf = "True";
			send_command(buf); 
			buf = "print(\"\")){print(\"\")";
			send_command(buf);
			
			switch_command();
			while (COMMAND_COUNTER < counter){
				send_command(command);					
				switch_command();		
			}; 

			buf = (("if (NOT(" + cond) + ")){print(\"\"");
			send_command(buf); 
			send_command("break");
			buf = "}"; 
			send_command(buf);
			send_command(buf); 
		}else{
			send_command(command); 
		};
	};	

	finish();
};

main();
