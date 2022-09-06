#import "stdlib/core.b"

string root_source;

void init(){
	root_source = "benv/prep_if_program.b";
	SET_SOURCE(root_source);
	SET_DEST("benv/while_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
	DEL_DEST(root_source); 
};

bool is_while(string command){
	stack s;
	string buf;

	buf = "while(";
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

	init();
	switch_command();

	for (print(""); print(""); NOT(command == "end"); switch_command()){
		if (is_while(command)){
			buf = "for (print(\"\")";
			send_command(buf); 
			buf = "print(\"\")";
			send_command(buf);
			send_command(get_cond(command)); 
			buf = "print(\"\")){print(\"\")";
			send_command(buf); 
		}else{
			send_command(command); 
		};
	};	

	finish();
};

main();
