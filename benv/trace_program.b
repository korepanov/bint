$filestdlib/files.b$;
$trace1$;
intopen_file(stringfile_path,stringmode){intres;
$trace2$;
intm;
$trace3$;
if(mode=="read"){m=0;
$trace4$;
}elseif(mode=="write"){m=1;
$trace5$;
}elseif(mode=="append"){m=2;
$trace6$;
}else{error="no such mode for opening file";
$trace7$;
if($toPanic){print((error+"\n"));
$trace8$;
exit(1);
$trace9$;
};
$trace10$;
returnres;
$trace11$;
};
$trace12$;
res=$open_f(file_path,m);
$trace13$;
returnres;
$trace14$;
};
$trace15$;
stringread_file(intdescriptor_number,intsize){intbytes;
$trace16$;
strings;
$trace17$;
bytes=$read_f(descriptor_number,size,s);
$trace18$;
returns;
$trace19$;
};
$trace20$;
voidwrite_file(intdescriptor_number,strings){intbytes;
$trace21$;
bytes=$write_f(descriptor_number,s);
$trace22$;
print(str(bytes));
$trace23$;
print("\n");
$trace24$;
if(bytes<0){error=("could not write to file with file descriptor number "+str(descriptor_number));
$trace25$;
if($toPanic){print((error+"\n"));
$trace26$;
exit(1);
$trace27$;
};
$trace28$;
};
$trace29$;
};
$trace30$;
voidclose_file(intdescriptor_number){bool$toPanicOld;
$trace31$;
$toPanicOld=$toPanic;
$trace32$;
$close_f(descriptor_number);
$trace33$;
try{strings;
$trace34$;
intbytes;
$trace35$;
intsize;
$trace36$;
size=1;
$trace37$;
bytes=$read_f(descriptor_number,size,s);
$trace38$;
};
$trace39$;
if(error==""){error=("could not close file with file descriptor number "+str(descriptor_number));
$trace40$;
if($toPanicOld){print((error+"\n"));
$trace41$;
exit(1);
$trace42$;
};
$trace43$;
};
$trace44$;
error="";
$trace45$;
};
$trace46$;
voiddel_file(stringfile_path){if(exists(file_path)){$del_f(file_path);
$trace47$;
}else{error=("no such file: "+file_path);
$trace48$;
if($toPanic){print((error+"\n"));
$trace49$;
exit(1);
$trace50$;
};
$trace51$;
};
$trace52$;
};
$fileprogram.b$;
$trace1$;
voidmain(){del_file("/home/slava/Go/bint/test.txt");
$trace2$;
print("");
$trace3$;
print("");
$trace4$;
print("");
$trace5$;
print("");
$trace6$;
};
$trace7$;
main();
