void main(){
	int number2;
	int number3;
	number2 = 2;
 	number3 = 3;
	SET_SOURCE((("program" + str(number2)) + ".b"));
	SET_DEST((("program" + str(number3)) + ".b"));
	string command;
	next_command(command);
	send_command(command);
	UNSET_SOURCE();
	UNSET_DEST(); 
};

main();
