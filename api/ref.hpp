// Copyright 2022 The Jule Programming Language.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

#ifndef __JULEC_REF_HPP
#define __JULEC_REF_HPP

constexpr signed int __JULEC_REFERENCE_DELTA{ 1 };

// Wrapper structure for raw pointer of JuleC.
// This structure is the used by Jule references for reference-counting
// and memory management.
template<typename T>
struct ref_jt;

template<typename T>
struct ref_jt {
    T *_alloc{ nil };
    mutable uint_jt *_ref{ nil };

    ref_jt<T>(void) noexcept: ref_jt<T>(T()) {}

    ref_jt<T>(std::nullptr_t) noexcept {}

    ref_jt<T>(T *_Ptr, uint_jt *_Ref) noexcept {
        this->_alloc = _Ptr;
        this->_ref = _Ref;
    }

    ref_jt<T>(T *_Ptr) noexcept {
        this->_ref = ( new( std::nothrow ) uint_jt );
        if (!this->_ref)
        { JULEC_ID(panic)( __JULEC_ERROR_MEMORY_ALLOCATION_FAILED ); }
        *this->_ref = 1;
        this->_alloc = _Ptr;
    }

    ref_jt<T>(const T &_Instance) noexcept {
        this->_alloc = ( new( std::nothrow ) T );
        if (!this->_alloc)
        { JULEC_ID(panic)( __JULEC_ERROR_MEMORY_ALLOCATION_FAILED ); }
        this->_ref = ( new( std::nothrow ) uint_jt );
        if (!this->_ref)
        { JULEC_ID(panic)( __JULEC_ERROR_MEMORY_ALLOCATION_FAILED ); }
        *this->_ref = __JULEC_REFERENCE_DELTA;
        *this->_alloc = _Instance;
    }

    ref_jt<T>(const ref_jt<T> &_Ref) noexcept
    { this->operator=( _Ref ); }

    ~ref_jt<T>(void) noexcept
    { this->__drop(); }

    inline int_jt __drop_ref(void) const noexcept
    { return ( __julec_atomic_add ( this->_ref, -__JULEC_REFERENCE_DELTA ) ); }

    inline int_jt __add_ref(void) const noexcept
    { return ( __julec_atomic_add ( this->_ref, __JULEC_REFERENCE_DELTA ) ); }

    inline uint_jt __get_ref_n(void) const noexcept
    { return ( __julec_atomic_load ( this->_ref ) ); }

    void __drop(void) noexcept {
        if (!this->_ref)
        { return; }
        if ( ( this->__drop_ref() ) != __JULEC_REFERENCE_DELTA )
        { return; }
        delete this->_ref;
        this->_ref = nil;
        delete this->_alloc;
        this->_alloc = nil;
    }

    inline T *operator->(void) noexcept
    { return ( this->_alloc ); }

    inline operator T(void) const noexcept
    { return ( *this->_alloc ); }

    inline operator T&(void) noexcept
    { return ( *this->_alloc ); }

    void operator=(const ref_jt<T> &_Ref) noexcept {
        this->__drop();
        if (_Ref._ref)
        { _Ref.__add_ref(); }
        this->_ref = _Ref._ref;
        this->_alloc = _Ref._alloc;
    }

    inline void operator=(const T &_Val) const noexcept
    { ( *this->_alloc ) = ( _Val ); }

    inline bool operator==(const ref_jt<T> &_Ref) const noexcept
    { return ( ( *this->_alloc ) == ( *_Ref._alloc ) ); }

    inline bool operator!=(const ref_jt<T> &_Ref) const noexcept
    { return ( !this->operator==( _Ref ) ); }

    friend inline
    std::ostream &operator<<(std::ostream &_Stream,
                             const ref_jt<T> &_Ref) noexcept
    { return ( _Stream << _Ref.operator T() ); }
};

#endif // #ifndef __JULEC_REF_HPP
