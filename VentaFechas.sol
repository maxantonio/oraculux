pragma solidity ^0.4.18;

interface IMyContract {
    function balanceOf(address owner) external returns (uint256); //conocer el balance
    function transfer(address to, uint256 amount) external returns (bool);//
}

contract DXTSale {
    IMyContract public tokenContract;  // direccion del token a venderse
    uint256 public price;              // el precio inicial
    uint256 public price1;              // precio despues de la primera fecha
    uint256 public price2;              // precio despues dela 2da fecha
    uint256 public price3;              // precio final

    address owner;
    uint256 public tokensSold;
    uint public start;
    uint public MES1;
    uint public MES2;
    uint public MES3;


    event Sold(address buyer, uint256 amount);

    function DXTSale(IMyContract _tokenContract, uint256 precio, uint256 precio1, uint256 precio2, uint256 precio3) public {
        owner = msg.sender;
        tokenContract = _tokenContract;
        price = precio * 1 finney;
        //0.1 finney;
        price1 = precio1 * 1 finney;
        price2 = precio2 * 1 finney;
        price3 = precio3 * 1 finney;
        start = now;
        MES1 = start + (DIA_SEG * 1);
        //tiempo mes 1 deberia ser 31 seteable despues en setMeses
        MES2 = start + (DIA_SEG * 2);
        //tiempo mes 2 deberia ser 61
        MES3 = start + (DIA_SEG * 3);
        //tiempo mes 3 deberia ser 92
    }

    function setMeses(uint timesTamp1, uint timesTamp2, uint timesTamp3){
        MES1 = timesTamp1;
        MES2 = timesTamp2;
        MES3 = timesTamp3;
    }

    function safeMultiply(uint256 a, uint256 b) internal pure returns (uint256) {
        if (a == 0) {
            return 0;
        } else {
            uint256 c = a * b;
            assert(c / a == b);
            return c;
        }
    }

    function() public payable {
        uint256 precio = getPrice();
        if (msg.value > 0) {
            uint numberOfTokens = msg.value / precio;
            uint256 scaledAmount = safeMultiply(numberOfTokens,
                uint256(10) ** 2);
            require(tokenContract.balanceOf(this) >= scaledAmount);
            emit Sold(msg.sender, numberOfTokens);
            tokensSold += numberOfTokens;
            require(tokenContract.transfer(msg.sender, scaledAmount));
        }
    }

    function getPrice() public returns (uint256 precio){
        precio = price;
        if (now > MES1) {
            precio = price1;
        }
        if (now > MES2) {
            precio = price2;
        }
        if (now > MES3) {
            precio = price3;
        }
    }

    function buyTokens() public payable {
        uint numberOfTokens = msg.value / price;
        uint256 scaledAmount = safeMultiply(numberOfTokens,
            uint256(10) ** 2);
        require(tokenContract.balanceOf(this) >= scaledAmount);
        emit Sold(msg.sender, numberOfTokens);
        tokensSold += numberOfTokens;
        require(tokenContract.transfer(msg.sender, scaledAmount));
    }

    function endSale() public {
        require(msg.sender == owner);
        // Send unsold tokens to the owner.
        require(tokenContract.transfer(owner, tokenContract.balanceOf(this)));
        msg.sender.transfer(address(this).balance);
    }
    //FUNCIONES PARA EL TRATO CON FECHA
    uint constant DIA_SEG = 86400;
    uint constant ANO_SEG = 31536000;
    uint constant ANO_BIS_SEG = 31622400;

    uint constant HORA_SEG = 3600;
    uint constant MIN_SEG = 60;


    uint16 constant PRIMER_ANO = 1970;

    function esBisiesto(uint16 ano) constant returns (bool) {
        if (ano % 4 != 0) {
            return false;
        }
        if (ano % 100 != 0) {
            return true;
        }
        if (ano % 400 != 0) {
            return false;
        }
        return true;
    }

    function toTimestamp(uint16 ano, uint8 mes, uint8 dia, uint8 hora, uint8 minuto, uint8 segundo) constant returns (uint timestamp) {
        uint16 i;

        // ano
        for (i = PRIMER_ANO; i < ano; i++) {
            if (esBisiesto(i)) {
                timestamp += ANO_BIS_SEG;
            }
            else {
                timestamp += ANO_SEG;
            }
        }

        // mes
        uint8[12] diaMes;
        diaMes[0] = 31;
        if (esBisiesto(ano)) {
            diaMes[1] = 29;
        }
        else {
            diaMes[1] = 28;
        }
        diaMes[2] = 31;
        diaMes[3] = 30;
        diaMes[4] = 31;
        diaMes[5] = 30;
        diaMes[6] = 31;
        diaMes[7] = 31;
        diaMes[8] = 30;
        diaMes[9] = 31;
        diaMes[10] = 30;
        diaMes[11] = 31;

        for (i = 1; i < mes; i++) {
            timestamp += DIA_SEG * diaMes[i - 1];
        }

        // dia
        timestamp += DIA_SEG * (dia - 1);

        // hora
        timestamp += HORA_SEG * (hora);

        // minuto
        timestamp += MIN_SEG * (minuto);

        // segundo
        timestamp += segundo;

        return timestamp;
    }
}