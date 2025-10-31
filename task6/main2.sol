// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title ERC-721 Non-Fungible Token Standard
/// @dev See https://eips.ethereum.org/EIPS/eip-721
///  Note: the ERC-165 identifier for this interface is 0x80ac58cd.
interface ERC721 {
    event Transfer(
        address indexed _from,
        address indexed _to,
        uint256 indexed _tokenId
    );
    event Approval(
        address indexed _owner,
        address indexed _approved,
        uint256 indexed _tokenId
    );
    event ApprovalForAll(
        address indexed _owner,
        address indexed _operator,
        bool _approved
    );

    function balanceOf(address _owner) external view returns (uint256);
    function ownerOf(uint256 _tokenId) external view returns (address);
    function safeTransferFrom(
        address _from,
        address _to,
        uint256 _tokenId,
        bytes memory data
    ) external payable;
    function safeTransferFrom(
        address _from,
        address _to,
        uint256 _tokenId
    ) external payable;
    function transferFrom(
        address _from,
        address _to,
        uint256 _tokenId
    ) external payable;
    function approve(address _approved, uint256 _tokenId) external payable;
    function setApprovalForAll(address _operator, bool _approved) external;
    function getApproved(uint256 _tokenId) external view returns (address);
    function isApprovedForAll(
        address _owner,
        address _operator
    ) external view returns (bool);
}

/// @title ERC-721 Non-Fungible Token Standard, optional metadata extension
/// @dev See https://eips.ethereum.org/EIPS/eip-721
///  Note: the ERC-165 identifier for this interface is 0x5b5e139f.
/* is ERC721 */ interface ERC721Metadata {
    function name() external view returns (string memory _name);
    function symbol() external view returns (string memory _symbol);

    /// @notice A distinct Uniform Resource Identifier (URI) for a given asset.
    /// @dev Throws if `_tokenId` is not a valid NFT. URIs are defined in RFC
    ///  3986. The URI may point to a JSON file that conforms to the "ERC721
    ///  Metadata JSON Schema".
    function tokenURI(uint256 _tokenId) external view returns (string memory);
}

// 在 NFT 转账时，检查接收方是否能安全接收 ERC-721 代币。
/// @dev Note: the ERC-165 identifier for this interface is 0x150b7a02.
interface ERC721TokenReceiver {
    function onERC721Received(
        address _operator,
        address _from,
        uint256 _tokenId,
        bytes memory _data
    ) external returns (bytes4);
}

// 检查实现了哪些接口。
interface ERC165 {
    function supportsInterface(bytes4 interfaceID) external view returns (bool);
}

contract MyNFT is ERC721, ERC721Metadata, ERC165 {
    string private _name = "My ERC721 NFT";
    string private _symbol = "MNFT";

    // 存储每个地址有几个 NFT
    mapping(address => uint256) private _balances;

    // 存储 Token 对应的 owner address
    mapping(uint256 => address) private _owners;
    // 存储 Token 对应的 authorized operator
    mapping(address => mapping(address => bool)) private _operators;
    // 存储 Token 对应的 approved address
    mapping(uint256 => address) private _approvers;

    // 保存 tokenID => CID
    mapping(uint256 => string) private _tokenURIs;

    constructor(string memory name_, string memory symbol_) {
        _name = name_;
        _symbol = symbol_;
    }

    function name() external view returns (string memory) {
        return _name;
    }

    function symbol() external view returns (string memory) {
        return _symbol;
    }

    function tokenURI(uint256 tokenId) external view returns (string memory) {
        address owner = _owners[tokenId];
        require(owner != address(0), "Invalid token.");

        return
            string(
                abi.encodePacked("https://ipfs.io/ipfs/", _tokenURIs[tokenId])
            );
    }

    function balanceOf(address owner) external view returns (uint256) {
        require(owner != address(0), "Invalid address.");

        return _balances[owner];
    }

    function ownerOf(uint256 tokenId) public view returns (address) {
        address owner = _owners[tokenId];
        require(owner != address(0), "Invalid token.");
        return owner;
    }

    // 授权 operator
    function setApprovalForAll(address operator, bool approved) external {
        require(operator != address(0), "Invalid address.");

        _operators[msg.sender][operator] = approved;
        emit ApprovalForAll(msg.sender, operator, approved);
    }

    function isApprovedForAll(
        address owner,
        address operator
    ) public view returns (bool) {
        require(
            owner != address(0) && operator != address(0),
            "Invalid address."
        );

        return _operators[owner][operator];
    }

    // 授权 approver
    function approve(address approver, uint256 tokenId) external payable {
        require(approver != address(0), "Invalid address.");

        address owner = ownerOf(tokenId);
        if (msg.sender != owner && !isApprovedForAll(owner, msg.sender)) {
            revert();
        }

        _approvers[tokenId] = approver;
        emit Approval(owner, approver, tokenId);
    }

    function getApproved(uint256 tokenId) public view returns (address) {
        return _approvers[tokenId];
    }

    function safeTransferFrom(
        address from,
        address to,
        uint256 tokenId,
        bytes memory data
    ) public payable {
        transferFrom(from, to, tokenId);

        // checks if `_to` is a smart contract (code size > 0)
        if (to.code.length > 0) {
            try
                ERC721TokenReceiver(to).onERC721Received(
                    msg.sender,
                    from,
                    tokenId,
                    data
                )
            returns (bytes4 response) {
                if (response != ERC721TokenReceiver.onERC721Received.selector) {
                    revert();
                }
            } catch (bytes memory) {
                revert();
            }
        }
    }

    function safeTransferFrom(
        address from,
        address to,
        uint256 tokenId
    ) external payable {
        safeTransferFrom(from, to, tokenId, "");
    }

    /// @dev
    ///     Throws unless `msg.sender` is the current owner, an authorized operator, or the approved address for this NFT.
    ///     Throws if `_from` is not the current owner.
    ///     Throws if `_to` is the zero address.
    ///     Throws if `_tokenId` is not a valid NFT.
    function transferFrom(
        address from,
        address to,
        uint256 tokenId
    ) public payable {
        require(from != address(0) && to != address(0), "Invalid address.");

        address owner = ownerOf(tokenId);
        // 三种身份可以执行
        if (
            msg.sender != owner &&
            !isApprovedForAll(owner, msg.sender) &&
            msg.sender != getApproved(tokenId)
        ) {
            revert();
        }

        // 只能是从所有者账号转出
        if (from != owner) {
            revert();
        }

        // 更换所有人
        _owners[tokenId] = to;

        // 减少计数
        _balances[from] -= 1;
        // 清除原来账户设置的授权地址。
        delete _approvers[tokenId];

        // 增加计数
        _balances[to] += 1;

        emit Transfer(from, to, tokenId);
    }

    function mint(address to, uint256 tokenId, string memory cid) public {
        require(to != address(0), "Invalid address.");

        address owner = _owners[tokenId];
        require(owner == address(0), "This tokenId has a owner.");

        _owners[tokenId] = to;
        // delete _approvers[tokenId]; 不清空了。
        _balances[to] += 1;

        _tokenURIs[tokenId] = cid;

        emit Transfer(address(0), to, tokenId);
    }

    function burn(uint256 tokenId) public {
        address owner = _owners[tokenId];
        require(owner != address(0), "This tokenId is invalid.");
        // 三种身份可以执行
        if (
            msg.sender != owner &&
            !isApprovedForAll(owner, msg.sender) &&
            msg.sender != getApproved(tokenId)
        ) {
            revert();
        }

        delete _owners[tokenId];
        delete _approvers[tokenId];
        _balances[owner] -= 1;

        emit Transfer(owner, address(0), tokenId);
    }

    function supportsInterface(
        bytes4 interfaceID
    ) external pure returns (bool) {
        return
            interfaceID == 0x80ac58cd || // ERC721
            interfaceID == 0x5b5e139f; //ERC721Metadata
    }
}
