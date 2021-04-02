
// var operand=[];
// var ope=[];
  
var user={
	number:[],
	operator:[]
}
function getHistory()
{
	return document.getElementById("history-value").innerText;
}

function getOutput()
{
	return document.getElementById("output-value").innerText;
}

function printOutput(num)
{
	if(num=="")
	{

		 document.getElementById("output-value").innerText=num;
	}
	else
	{
		 document.getElementById("output-value").innerText=getFormattedNumber(num);
	}
}
function printHistory(num)
{
	 document.getElementById("history-value").innerText=num;
}

function getFormattedNumber(num)
{
	var n=Number(num);
	var value=n.toLocaleString("en");
	return value;
}

function reverseFormattedNumber(num)
{
   return Number(num.replace(/,/g,''));

}

var operator=document.getElementsByClassName("operator")
for (var i = 0; i <operator.length; i++) {
operator[i].addEventListener('click',function()
{ 
	if(this.id=="clear")
	{
		printOutput("");
		printHistory("");
	}
	else if(this.id=="backspace")
	{
		var output=reverseFormattedNumber(getOutput()).toString();
		if (output)
		{
			output=output.substr(0,output.length-1);
			printOutput(output);
		}
	}
	else{
		var output=getOutput();
		
		var history=getHistory();
		if(output!="")
		{
			output=reverseFormattedNumber(output);
			user.number.push(String(output))
			console.log(user.number)
			history=history+output;
			if(this.id=="=")
			{	console.log(JSON.stringify(user))
				fetch('http://localhost:9000/calc', {
				method: 'POST',
				body: JSON.stringify(user)
				})
				.then(response => response.json())
				.then(result =>printOutput(JSON.stringify(result)))
				// var result=eval(history);
				printHistory("");
				user.number=[]
				user.operator=[]
			}
			else
			{
				user.operator.push(this.id)
				console.log(user.operator)
				history=history+this.id;
				printHistory(history);
				printOutput("");
			}
			
		}
		
	}
}
)
}

//Getting Numbers
var output
var number=document.getElementsByClassName('number');
for (var i = 0; i < number.length; i++) {
	number[i].addEventListener('click',function()
	{
		var output=reverseFormattedNumber(getOutput());
     if(output!=NaN)
     {
		output=output+this.id;
     	printOutput(output);
		 	
     }	
	}
	
	)	
	
}
